package main

import (
	assets "app"
	types "common"
	"fmt"
	"os"
)


const template_docx = "template.docx"
const template_xlsx = "template.xlsx"
const template_pptx = "template.pptx"
const template_txt = "template.txt"

const tmp_prefix = "gravtmp"
const tmp_docx = "gravtmp_*.docx"
const tmp_xlsx = "gravtmp_*.xlsx"
const tmp_pptx = "gravtmp_*.pptx"
const tmp_txt = "gravtmp_*.txt"

type IRunner interface {
	SetupEnv() error
	RestoreEnv() error
	NewTemplate(types.AppType) (string, error)
	// waits until app is finished runninig
	OpenApp(typ types.AppType, file_path string) error
	FocusOrOpenApp(typ types.AppType, file_path string) error
	FocusOpenApp() error
	KillApp() error
	ListAllProcess() (map[uint32]string, error)
}

func (self *Runner) NewTemplate(typ types.AppType) (string, error) {
	var tmp string
	var template string
	switch typ {
	case types.TXT:
		tmp = tmp_txt
		template = template_txt
	case types.DOCX:
		tmp = tmp_docx
		template = template_docx
	case types.PPTX:
		tmp = tmp_pptx
		template = template_pptx
	case types.XLSX:
		tmp = tmp_xlsx
		template = template_xlsx
	default:
		return "", fmt.Errorf("unknown app type %d", typ)
	}

	file, err := os.CreateTemp("", tmp)
	if err != nil {
		return "", err
	}
	file.Close()

	dest := file.Name()
	err = self.newTemplate(template, dest)
	if err != nil {
		return "", err
	}
	return dest, nil
}

func (self *Runner) newTemplate(name string, dest string) error {
	// NOTE: non os specific path separaters
	path := fmt.Sprintf("templates/%s", name)
	contents, err := assets.Templates.ReadFile(path)
	if err != nil {
		return err
	}
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(contents)
	return err
}





// func (self *Runner) run(name string, args ...string) error {
// 	cmd := exec.Command(name, args...)
// 	err := cmd.Run()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return err
// }
