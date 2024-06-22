package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content, nil
}

func writeFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func interpretTagmake(text string) (string, string, error) {
	lines := strings.Split(text, "\n")

	var out []string

	namespace := ""
	isTag := false

	outpath := ""
	replaces := 0

	collecting := false
	for ln, l := range lines {
		l = strings.TrimSpace(strings.Split(l, ",")[0])
		if l == "" {
			continue
		} else if collecting {
			if l == ";" {
				collecting = false
				namespace = ""
				isTag = false
				continue
			}

			if isTag {
				out = append(out, "        \"#"+namespace+l+"\"")
			} else {
				out = append(out, "        \""+namespace+l+"\"")
			}
		} else {
			if l[len(l)-1] == ':' {
				args := strings.Split(l, " ")

				name := strings.TrimSpace(args[len(args)-1])
				if name[len(name)-1] != ':' {
					return "", "", fmt.Errorf("error on line %d: arguments are only allowed on the left side of the namespace name", ln)
				}
				namespace = name
				if namespace[len(namespace)-1] != ':' {
					namespace += ":"
				}

				if len(args) > 1 {
					args = args[:len(args)-1]

					typeAlreadyGiven := false
					for _, a := range args {
						switch strings.TrimSpace(a) {
						case "":
							continue
						case "item":
							if typeAlreadyGiven {
								return "", "", fmt.Errorf("error on line %d: cannot assign type '%s', type was already given", ln, strings.TrimSpace(a))
							}
							typeAlreadyGiven = true
						case "tag":
							if typeAlreadyGiven {
								return "", "", fmt.Errorf("error on line %d: cannot assign type '%s', type was already given", ln, strings.TrimSpace(a))
							}
							typeAlreadyGiven = true
							isTag = true
						default:
							return "", "", fmt.Errorf("error on line %d: invalid argument '%s'", ln, a)
						}
					}
				}

				collecting = true
			} else if len(l) >= 5 {
				if l[:5] == "!out " {
					if outpath != "" {
						return "", "", fmt.Errorf("error on line %d: output path cannot be redeclared", ln)
					} else if strings.TrimSpace(l[len(l)-5:]) == "" {
						return "", "", fmt.Errorf("error on line %d: output path cannot be empty", ln)
					} else {
						outpath = l[5:]
					}
				} else {
					return "", "", fmt.Errorf("error on line %d: invalid line '%s'", ln, l)
				}
			} else if len(l) >= 10 {
				if l[:10] == "!replaces " {
					if replaces != 0 {
						return "", "", fmt.Errorf("error on line %d: replaces cannot be redeclared", ln)
					} else {
						switch strings.ToLower(strings.TrimSpace(l[len(l)-10:])) {
						case "true":
							replaces = 2
						case "false":
							replaces = 1
						default:
							return "", "", fmt.Errorf("error on line %d: replaces input must be 'true' or 'false'", ln)
						}
					}
				} else {
					return "", "", fmt.Errorf("error on line %d: invalid line '%s'", ln, l)
				}
			} else {
				return "", "", fmt.Errorf("error on line %d: invalid line '%s'", ln, l)
			}
		}
	}

	if outpath == "" {
		return "", "", fmt.Errorf("error: TagMake files must have a '!out [output path]' in them")
	}

	replacesFinal := false
	if replaces == 2 {
		replacesFinal = true
	}
	return fmt.Sprintf("{\n    \"replace\": %v,\n    \"values\": [\n", replacesFinal) + strings.Join(out, ",\n") + "\n    ]\n}", path.Clean(outpath), nil
}
