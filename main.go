// This utility will concurrently walk through the specified dir(s)
// and convert the images found to webp, it will create a new webp file,
// that has the name of the old file + and appended string ".webp", and you can optionally 
// clean the old file with the "--clean" switch.
package main

import "os"
import "log"
import "flag"
import "sync"
import "regexp"
import "strings"
import "path/filepath"
import "gopkg.in/h2non/bimg.v1"

var (
	DIR   = flag.String("dir", "/var/www/uploads,/var/www/assets", "director(y|ies) to search in")
	CLEAN = flag.Bool("clean", false, "whether to delete the old images or not")
	EXT   = flag.String("ext", "png,jpg,jpeg", "extensions to process")
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	re := regexp.MustCompile("(?i).(" + strings.Replace(*EXT, ",", "|", -1) + ")$")
	for _, searchDir := range strings.Split(*DIR, ",") {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()
			filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					log.Println(err)
					return nil
				}
				if !re.MatchString(path) {
					return nil
				}
				wg.Add(1)
				go func() {
					defer wg.Done()
					convert(path, re.ReplaceAllString(path, ".webp"), *CLEAN)
				}()
				return nil
			})
		}(searchDir)
	}
	wg.Wait()
}

// convert the specified $in path to webp, and 
// save it as $out, and optionally $clean if needed.
func convert(in, out string, clean bool) error {
	log.Println("[processing]", in)
	buffer, err := bimg.Read(in)
	if err != nil {
		log.Println(err)
		return err
	}
	newImage, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		log.Println(err)
		return err
	}
	err = bimg.Write(out, newImage)
	if err != nil {
		log.Println(err)
	}
	if clean {
		os.Remove(in)
	}
	return nil
}
