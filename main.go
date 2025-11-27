package main

import (
	"encoding/binary"
	"image/color"
	"net"

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

	packetCh := make(chan Pixel, 1024)

	img := NewImage()

	go func() {
		buf := make([]byte, 1024)
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
		for p := range packetCh {
			x := int(p.index % 256)
			y := int(p.index / 256)
			img.pixels[y][x] = p.color
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

func sendReSignal(addr *net.UDPAddr) {

}
