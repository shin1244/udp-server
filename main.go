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

	drawCh := make(chan Pixel)
	checkCh := make(chan Pixel)

	img := NewImage()

	go readUDPPackets(conn, drawCh, checkCh)
	go requestMissingPixels(conn, checkCh)
	go drawPixels(img, drawCh)

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

func readUDPPackets(conn *net.UDPConn, drawCh chan<- Pixel, checkCh chan<- Pixel) {
	buf := make([]byte, 128)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		p := parsePixel(buf[:n], addr)
		drawCh <- p
		checkCh <- p
	}
}

func requestMissingPixels(conn *net.UDPConn, checkCh <-chan Pixel) {
	received := make([]bool, 65536)
	reqBuf := make([]byte, 2)

	var ticker *time.Ticker
	var lastAddr *net.UDPAddr

	for {
		select {
		case p := <-checkCh:
			received[p.index] = true

			if p.index == 65535 {
				if ticker != nil {
					ticker.Stop()
				}
				ticker = time.NewTicker(1000 * time.Millisecond)
				lastAddr = p.addr
			}

		case <-ticker.C:
			missingCount := 0

			for i := 0; i < 65536; i++ {
				if !received[i] {
					binary.BigEndian.PutUint16(reqBuf, uint16(i))
					conn.WriteToUDP(reqBuf, lastAddr)
					missingCount++
				}
			}
			if missingCount == 0 {
				fmt.Println("All pixels received!")
				ticker.Stop()
				ticker = nil
			} else {
				fmt.Printf("Requested %d missing packets.\n", missingCount)
			}
		}
	}
}

func drawPixels(img *Image, drawCh <-chan Pixel) {
	for p := range drawCh {
		x := int(p.index % 256)
		y := int(p.index / 256)
		img.pixels[y][x] = p.color
	}
}
