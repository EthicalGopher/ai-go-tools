package tts

import (
	"flag"
	"fmt"
	"os"


	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)



func TTS(input,filepath,folderpath string) (string,error) {
	textPtr := flag.String("text", input, "Text to convert to speech")
	langPtr := flag.String("lang", voices.English, "Language code (e.g., en, fr, es, de, etc.)")
	outputDirPtr := flag.String("outdir", folderpath, "Directory to save the audio file")
	outputFilePtr := flag.String("outfile", filepath, "Name of the output file (without extension)")
	flag.Parse()


	
	path,_:= os.Stat(folderpath);
	if path!=nil{
		err:=os.Remove(folderpath+`\`+filepath+`.mp3`)
		if err!=nil{
			return "",err
			}
		
	}
	// Ensure output directory exists
	if _, err := os.Stat(*outputDirPtr); os.IsNotExist(err) {
		err := os.MkdirAll(*outputDirPtr, 0755)
		if err != nil {
			return "",err
		}
	}


	// Initialize the TTS system
	speech := htgotts.Speech{
		Folder:   *outputDirPtr,
		Language: *langPtr,
		Handler:  &handlers.Native{},
	}

	// Convert text to speech
	fmt.Printf("Converting text to speech: '%s'\n", *textPtr)
	_, err := speech.CreateSpeechFile(*textPtr, *outputFilePtr)
	if err != nil {
		return "",err
	}


	return "success",nil
}
