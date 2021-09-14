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

var (
	Version = ""

	AppVersions = map[string]string{
		"v012c": "2mhihr",
		"v012b": "clrt9k",
		"v012a": "l2a2vu",
		"v011e": "clrt9k",
		"v011d": "1uyyh4",
		"v011c": "bzqbww",
		"v011b": "s0xirg",
		"v011a": "khhfux",
	}

	StreamableURL = "https://streamable.com/e"

	Slides = []string{"PREV VERSION", "NEXT VERSION"}
)

func (m *Model) UpdateVersion(newVersion string) {
	m.Version = newVersion
}

func (m *Model) StreamSrcVersion(newVersion string) string {
	return fmt.Sprintf("%s/%s", m.Streamable, m.Versions[m.Version])
}

func (m *Model) SwapVersion() {
	// log.Println("version", m.Version, m.VersionKeys, m.VersionIdx)
	m.Version = m.VersionKeys[m.VersionIdx]
	slideIdxs := [2]int{}
	slideIdxs[0] = (m.VersionsCount + (m.VersionIdx-m.VersionsCount)%m.VersionsCount) - 1
	slideIdxs[1] = (m.VersionIdx + 1) % m.VersionsCount
	tmpSlides := []string{
		m.VersionKeys[slideIdxs[0]],
		m.VersionKeys[slideIdxs[1]],
	}
	m.Slides = tmpSlides
	// log.Printf("next: %s prev: %s", m.VersionKeys[slideIdxs[1]], m.VersionKeys[slideIdxs[0]])
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
			`<a style="color: #fff; text-decoration: none;" href="%s">%s — <strong>REVZIM</strong></a>`,
			vm.Data.Get("authorgit").String(), time.Now().Format("Mon Jan _2 2006"),
		)
	})

	// o = o.Mixin(js.M{
	// 	"vuetify": InitVuetify(),
	// })
	return o
}

func main() {
	// or as part of a structed object:
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}
	m.MountEl = "#app"
	m.Version = ""
	m.Versions = AppVersions
	m.VersionIdx = 0
	m.VersionKeys = make([]string, 0)
	m.Slides = Slides
	m.AuthorGit = "https://github.com/revzim"
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

	// vueApp := v.Mount("#app")

	js.Global.Set("app", v)
}
