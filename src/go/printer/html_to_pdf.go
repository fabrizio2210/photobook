package main

import (
  "bytes"
	"context"
  "fmt"
	"io/ioutil"
	"log"
  "os"
  "os/exec"
  "strings"
  "text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

  "Lib/filemanager"
  "Lib/models"
)

var htmlTemplate = `
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" lang="" xml:lang="">
<head>
<meta name="generator" content=
"HTML Tidy for HTML5 for Linux version 5.6.0" />
<title>invito-html.html</title>
<meta http-equiv="Content-Type" content=
"text/html; charset=utf-8" />
<style type="text/css">
.vignetta {width: min-content;}
p {margin: 0; padding: 0; text-align: center; white-space:break-spaces; overflow-wrap: anywhere}
.ft11{font-size:30px;font-family:BAAAAA+GlacialIndifference;color:#2c3e50;}
.ft16{font-size:34px;font-family:AAAAAA+GlacialIndifference;color:#2c3e50;}
</style>
</head>
<body bgcolor="#FFFFFF" vlink="blue" link="blue" style="width:297mm; height:208mm" >
<div id="page1-div" style="display:flex; justify-content:space-around; margin:0; border: 0;width:295mm;height:207mm;position:relative">

<img style="position:absolute;top:10mm;right:10mm;width:100mm;z-index:0" src="resources/right_leaf.png" alt="leaf" />
<img style="position:absolute;bottom:10mm;left:10mm;width:100mm;z-index:-1" src="resources/left_leaf.png" alt="leaf" />

{{ range $i, $p := .Photos }}

{{ if $p.Photo_id }}
{{ if eq $i 0 }} 
<div class="vignetta" style="">
{{ else }}
<div class="vignetta" style="align-self:flex-end">
{{ end }}

<img style="max-width:140mm;max-height:140mm" src="{{ $p.Location }}" alt="img0"/>
<p>{{ $p.Description }}</p>
<p>-- {{ $p.Author }} --</p>

</div>
{{ end }}
{{ end }}

</div>
</body>
</html>
`

type contentPage struct {
  Photos *[2]models.PhotoEvent
}

func commandToPrintSinglePDF(urlstr string, res *[]byte) chromedp.Tasks {
  // TODO compress images during print
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
                          WithPrintBackground(true).
                          WithPreferCSSPageSize(true).
                          WithPaperWidth(8.28).
                          WithPaperHeight(11.7).
                          WithLandscape(true).Do(ctx) // PrintToPDF through cdp implementation
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

func printSinglePDF(ctx context.Context, page *[2]models.PhotoEvent, intermediatePDF string ) {
  // Prepare working directory
  tempDir, err := ioutil.TempDir("/tmp/", "single-")
  if err != nil {
	  log.Fatal(err)
  }
  defer os.RemoveAll(tempDir)

  cwd, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  err = os.Symlink(cwd + "/resources", tempDir + "/resources")
  if err != nil {
    log.Fatal(err)
  }

  htmlFile := tempDir + "/page.html"
  content := contentPage{
    Photos: page,
  }
  t, err := template.New("pagina").Parse(htmlTemplate)
  var htmlContent bytes.Buffer
  err = t.Execute(&htmlContent, content)
  if err != nil {
    log.Fatal(err)
  }
  os.WriteFile(htmlFile, htmlContent.Bytes(), 0644)
	// Generate pdf
	var buf []byte
	if err := chromedp.Run(ctx, commandToPrintSinglePDF("file://" + htmlFile, &buf)); err != nil {
		log.Fatal(err)
	}
  log.Printf("Writing %s", intermediatePDF)
	if err := ioutil.WriteFile(intermediatePDF, buf, 0644); err != nil {
		log.Fatal(err)
	}
}


func printToPDF(outputFile string, layout []*[2]models.PhotoEvent) {
  opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("headless", true),
  )
  ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
  defer cancel()
  ctx, cancel = chromedp.NewContext(ctx)
  defer cancel()
   
  tempDir, err := ioutil.TempDir("/tmp/", "intermediate-")
  if err != nil {
	  log.Fatal(err)
  }
  defer os.RemoveAll(tempDir)

  var intermediateFiles []string
  _, err = os.Stat(filemanager.GetCoverLocation())
  if ! os.IsNotExist(err) {
    intermediateFiles = append (intermediateFiles, filemanager.GetCoverLocation())
  }
  for i, page := range layout {
    intermediateFile := fmt.Sprintf("%s/%04d.pdf", tempDir, i)
    printSinglePDF(ctx, page, intermediateFile)
    intermediateFiles = append(intermediateFiles, intermediateFile)
  }
  intermediateFiles = append(intermediateFiles, outputFile)
  log.Printf("Unifying in %s", outputFile)
  cmd := exec.Command("pdfunite", intermediateFiles...)
  var out strings.Builder
  cmd.Stdout = &out
  err = cmd.Run()
  if err != nil {
    log.Print(err)
  }
  log.Printf("Output: %s", out.String())
}
