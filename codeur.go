// codeur
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	TAB string = "    "
)

var gens = map[string] string {
	"[NOM]"  : "",
	"[STYLE]": "",
	"[CORPS]": "",
	"[ICONE]": "",
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func traiterLigne(ligne string) (rendu string) {
	ind := strings.Index(ligne, ":")
	if ind > -1 {
		motcle := trim(ligne[:ind])
		valeur := trim(ligne[ind+1:])
		switch motcle {
		case "nom":
			rendu = "<title>" + valeur + "</title>"
			gens["[NOM]"] = rendu
		case "icone":
			gens["[ICONE]"] = "<link rel=\"icon\" type=\"image/png\" href=\""+ valeur +"\" />"
		case "style":
			if strings.Contains(valeur, ",") {
				styles := strings.Split(valeur, ",")
				for i := 0; i<len(styles); i++ {
					if i > 0 {
						gens["[STYLE]"] += "\n        "
					}
					gens["[STYLE]"] += "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + trim(styles[i]) + ".css\" />"
				}
			} else {
				rendu = "<link type=\"text/css\" rel=\"stylesheet\" href=\"" + valeur + ".css\" />"
				gens["[STYLE]"] = rendu
			}

		}
	}
	return
}

func trim(text string) string {
	return strings.Trim(text, " ")
}

func contenu(chemin string) string {
	dat, err := ioutil.ReadFile(chemin)
	check(err)
	return string(dat)
}

func main() {

	args := os.Args
	var fichier string
	if len(args) > 1 {
		fichier = args[1]
	} else {
		return
	}

	file, err := os.Open(fichier)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ligne := traiterLigne(scanner.Text())
		if ligne != "" {
			fmt.Println(ligne)
		}
	}

	check(scanner.Err())

	template := contenu("template.html")
	sortie := strings.Replace(template, "[NOM]", gens["[NOM]"], 1)
	sortie = strings.Replace(sortie, "[STYLE]", gens["[STYLE]"], 1)
	sortie = strings.Replace(sortie, "[ICONE]", gens["[ICONE]"], 1)
	d1 := []byte(sortie)
	check(ioutil.WriteFile("page.html", d1, 0644))
	cmd := exec.Command("xdg-open", "page.html")
	check(cmd.Start())
}
