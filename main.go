package main

import (
	"flag"
	"fmt"
	"github.com/OGFris/MP3ToVideo/mjpeg"
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/dhowden/tag"
	"github.com/tcolgate/mp3"
	"image"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func main() {
	var in, out string

	flag.StringVar(&in, "i", "", "give input file")
	flag.StringVar(&out, "o", "", "give output file")
	flag.Parse()

	if in == "" {
		log.Fatalln("please give a valid input file")
	}

	if out == "" {
		log.Fatalln("please give a valid output file")
	}

	if _, err := os.Stat(out); err == nil {
		err = os.Remove(out)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(in); err != nil {
		log.Fatalln("couldn't find the input file")
	}

	f, err := os.Open(in)
	if err != nil {
		panic(err)
	}

	m, err := tag.ReadFrom(f)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	temp := fmt.Sprint(os.TempDir()+"/", rand.Float32(), ".avi")
	aw, err := mjpeg.New(temp, 1920, 1080, 1)
	if err != nil {
		panic(err)
	}

	d := mp3.NewDecoder(f)

	var (
		frame   mp3.Frame
		skipped = 0
		dur     float64
	)
	for {
		err := d.Decode(&frame, &skipped)
		if err != nil {
			break
		}

		dur = dur + frame.Duration().Seconds()

	}

	imgF, err := os.Create(temp + ".jpeg")
	if err != nil {
		panic(err)
	}

	_, err = imgF.Write(m.Picture().Data)
	if err != nil {
		panic(err)
	}

	err = imgF.Close()
	if err != nil {
		panic(err)
	}

	img, err := imgio.Open(temp + ".jpeg")
	if err != nil {
		panic(err)
	}

	cropped := transform.Crop(img, image.Rect(0, 100, 1920, 800))
	resized := transform.Resize(cropped, 1920, 1080, transform.Linear)
	result := blur.Gaussian(resized, 50.0)
	bright := adjust.Brightness(result, -0.7)

	if err = imgio.Save(temp+"final.jpeg", bright, imgio.JPEGEncoder(100)); err != nil {
		panic(err)
	}

	for i := 0; i < int(dur); i++ {
		imgf, err := os.Open(temp + "final.jpeg")
		if err != nil {
			panic(err)
		}
		imgBytes, err := ioutil.ReadAll(imgf)
		if err != nil {
			panic(err)
		}

		err = aw.AddFrame(imgBytes)
		if err != nil {
			panic(err)
		}
	}

	err = aw.Close()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("ffmpeg", "-i", temp, "-i", in, "-codec", "copy", "-shortest", temp+".avi")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd2 := exec.Command("ffmpeg", "-i", temp+".avi", out)
	if err := cmd2.Run(); err != nil {
		panic(err)
	}
}
