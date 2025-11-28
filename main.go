package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct {
	pixels [256][256]color.RGBA
}

type Pixel struct {
	index uint16
	color color.RGBA
	addr  *net.UDPAddr
}

func main() {
	laddr, _ := net.ResolveUDPAddr("udp", ":8080")
	conn, _ := net.ListenUDP("udp", laddr)
	defer conn.Close()

	packetCh := make(chan Pixel)
	checkCh := make(chan Pixel)

	img := NewImage()

	go func() {
		buf := make([]byte, 128)
		for {
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				continue
			}
			p := parsePixel(buf[:n], addr)
			packetCh <- p
		}
	}()

	go func() {
		// 1차원 배열로 체크하는 게 인덱스 계산하기 더 편함
		received := make([]bool, 65536)

		var verifyTimer <-chan time.Time
		var lastAddr *net.UDPAddr

		for {
			select {
			case p := <-checkCh:
				received[p.index] = true

				if p.index == 65535 {
					verifyTimer = time.Tick(1000 * time.Millisecond)
					lastAddr = p.addr
				}

			case <-verifyTimer:
				missingCount := 0

				reqBuf := make([]byte, 2)

				for i := 0; i < 65536; i++ {
					if !received[i] {
						binary.BigEndian.PutUint16(reqBuf, uint16(i))

						if lastAddr != nil {
							conn.WriteToUDP(reqBuf, lastAddr)
						}

						missingCount++

						if missingCount%100 == 0 {
							time.Sleep(time.Millisecond * 1)
						}
					}
				}

				if missingCount == 0 {
					fmt.Println("완벽함! 모든 픽셀 수신 완료.")
					verifyTimer = nil
				} else {
					fmt.Printf("총 %d개의 패킷 재요청 보냄.\n", missingCount)
				}

			}
		}
	}()

	go func() {
		for p := range packetCh {
			x := int(p.index % 256)
			y := int(p.index / 256)
			img.pixels[y][x] = p.color
			checkCh <- p
		}
	}()
	ebiten.SetWindowSize(256, 256)
	ebiten.SetWindowTitle("Image Viewer")

	ebiten.RunGame(img)
}

func NewImage() *Image {
	return &Image{
		pixels: [256][256]color.RGBA{},
	}
}

func parsePixel(data []byte, addr *net.UDPAddr) Pixel {
	return Pixel{
		index: binary.BigEndian.Uint16(data[0:2]),
		color: color.RGBA{R: data[2], G: data[3], B: data[4], A: data[5]},
		addr:  addr,
	}
}

func (i *Image) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 256, 256
}

func (i *Image) Update() error {
	return nil
}

func (i *Image) Draw(screen *ebiten.Image) {
	for y := range 256 {
		for x := range 256 {
			screen.Set(x, y, i.pixels[y][x])
		}
	}
}
