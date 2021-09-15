package main

// revzim <https://github.com/revzim>
// aztd - demo website

import (
	"fmt"
	"strings"
	"time"

	vue "github.com/revzim/gopherjs-vue"

	"github.com/gopherjs/gopherjs/js"
)

type (
	Model struct {
		*js.Object
		MountEl       string            `js:"mountEl"`
		Version       string            `js:"version"`
		Versions      map[string]string `js:"versions"`
		Slides        []string          `js:"slides"`
		VersionKeys   []string          `js:"versionKeys"`
		VersionsCount int               `js:"versionsCount"`
		VersionIdx    int               `js:"versionIdx"`
		Streamable    string            `js:"streamable"`
		AuthorGit     string            `js:"authorgit"`
		Vuetify       *js.Object        `js:"vuetify"`
	}
)

const (
	VueAppMountElement = "#app"
	AuthorGit          = "https://github.com/revzim"
	FooterFormatString = `<a style="color: #fff; text-decoration: none;" target="_blank" href="%s">%s — <strong>REVZIM</strong></a>`
	StreamableURL      = "https://streamable.com/e"
)

var (
	AppVersions = map[string]string{
		"v013b": "kgeu1w",
		"v013a": "xs6kid",
		"v012c": "2mhihr",
		"v012b": "clrt9k",
		"v012a": "l2a2vu",
		"v011e": "clrt9k",
		"v011d": "1uyyh4",
		"v011c": "bzqbww",
		"v011b": "s0xirg",
		"v011a": "khhfux",
	}

	Slides = []string{"PREV VERSION", "NEXT VERSION"}
)

func (m *Model) UpdateVersion(newVersion string) {
	m.Version = newVersion
}

func (m *Model) StreamSrcVersion(newVersion string) string {
	return fmt.Sprintf("%s/%s", m.Streamable, m.Versions[m.Version])
}

func (m *Model) SwapVersion() {
	m.Version = m.VersionKeys[m.VersionIdx]
	const arrLen = 2
	slideIdxs := [arrLen]int{}
	slideIdxs[0] = (m.VersionsCount + (m.VersionIdx-m.VersionsCount)%m.VersionsCount) - 1
	slideIdxs[1] = (m.VersionIdx + 1) % m.VersionsCount
	tmpSlides := make([]string, arrLen)
	for i := range slideIdxs {
		tmpSlides[i] = m.VersionKeys[slideIdxs[i]]
	}
	m.Slides = tmpSlides
}

func InitVuetify() *js.Object {
	Vuetify := js.Global.Get("Vuetify")

	vue.Use(Vuetify)

	vuetifyCFG := js.Global.Get("Object").New()
	vuetifyCFG.Set("theme", map[string]interface{}{
		"dark": true,
	})

	vtfy := Vuetify.New(vuetifyCFG)

	return vtfy
}

func InitVueOpts(m *Model) *vue.Option {
	o := vue.NewOption()

	o.SetDataWithMethods(m)

	o.AddComputed("versionLabel", func(vm *vue.ViewModel) interface{} {
		return strings.ToUpper(vm.Data.Get("version").String())
	})

	o.AddComputed("currentSlides", func(vm *vue.ViewModel) interface{} {
		return vm.Data.Get("slides")
	})

	o.AddComputed("footerHTML", func(vm *vue.ViewModel) interface{} {
		return fmt.Sprintf(
			FooterFormatString,
			vm.Data.Get("authorgit").String(),
			time.Now().Format("Mon Jan _2 2006"),
		)
	})

	o = o.Mixin(js.M{
		"vuetify": InitVuetify(),
	})

	return o
}

func main() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}
	m.MountEl = VueAppMountElement
	m.Version = ""
	m.Versions = AppVersions
	m.VersionIdx = 0
	m.VersionKeys = make([]string, 0)
	m.Slides = Slides
	m.AuthorGit = AuthorGit
	m.Streamable = StreamableURL
	for k := range m.Versions {
		m.VersionKeys = append(m.VersionKeys, k)
	}
	m.VersionsCount = len(m.VersionKeys)

	m.SwapVersion()

	m.Vuetify = InitVuetify()

	o := InitVueOpts(m)

	v := o.NewViewModel()

	v.Object.Set("vuetify", m.Vuetify)

	v.Mount(VueAppMountElement)
}
