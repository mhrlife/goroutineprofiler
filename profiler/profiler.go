package profiler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Profiler struct {
	Config
	lastMaxCount int
}

type Config struct {
	Sleep     time.Duration
	MaxWindow int // MB
	ServerURL string
}

func DefaultConfig() Config {
	return Config{
		Sleep:     time.Second * 20,
		MaxWindow: 20,
		ServerURL: "http://localhost:6060",
	}
}

func NewProfiler(c Config) *Profiler {
	return &Profiler{
		Config: c,
	}
}

func (p *Profiler) Run() {
	for {
		p.work()
		time.Sleep(p.Sleep)
	}
}

func (p *Profiler) work() {
	fileName := p.saveProfile()
	totalRoutines := p.extractTotalRoutines(fileName)
	if p.isInSaveRange(totalRoutines) {
		p.lastMaxCount = totalRoutines
		p.renameFile(fileName, totalRoutines)
	} else {
		p.deleteFile(fileName)
	}
}

func (p *Profiler) saveProfile() string {
	fileName := "sample-" + time.Now().Format(time.RFC3339) + "-COUNT.out"
	pprofUrl := p.pprofGoroutineURL()
	err := downloadFile(p.localFilePath(fileName), pprofUrl)
	if err != nil {
		fmt.Println("error while downloading ", err)
		return ""
	}
	return fileName
}

func (p *Profiler) extractTotalRoutines(fileName string) int {
	cmd := exec.Command("go", "tool", "pprof", "-text", p.localFilePath(fileName))
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error while creating", err)
		return 0
	}
	output := string(out)
	lines := strings.Split(output, "\n")
	firstLine := strings.Split(strings.TrimSpace(lines[5]), " ")
	totalRoutines := firstLine[0]
	t, _ := strconv.Atoi(totalRoutines)
	return t
}

func (p *Profiler) renameFile(fileName string, count int) {
	newFileName := strings.ReplaceAll(fileName, "COUNT", strconv.Itoa(count))
	os.Rename(p.localFilePath(fileName), p.localFilePath(newFileName))
}

func (p *Profiler) deleteFile(fileName string) {
	os.Remove(p.localFilePath(fileName))
}

func (p *Profiler) pprofGoroutineURL() string {
	return p.ServerURL + "/debug/pprof/goroutine"
}

func (p *Profiler) isInSaveRange(count int) bool {
	return count > p.lastMaxCount+p.MaxWindow
}

func (p *Profiler) localFilePath(fileName string) string {
	return "./samples/" + fileName
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
