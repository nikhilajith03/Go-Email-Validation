package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"text/template"
)

// Step 1: Email Format Validation
func isValidEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Step 2: MX Record Lookup (Domain Check)
func isValidDomain(domain string) bool {
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}
	return true
}

// Step 3: SMTP Check (Verifying the Email Exists)
func checkSMTP(email string) bool {
	domain := strings.Split(email, "@")[1]
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}

	server := mxRecords[0].Host
	client, err := smtp.Dial(fmt.Sprintf("%s:25", server))
	if err != nil {
		return false
	}
	defer client.Close()

	client.Hello("localhost")
	client.Mail("sender@example.com")
	err = client.Rcpt(email)
	if err != nil {
		return false
	}
	return true
}

// HTTP handler to serve the HTML page
func emailCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		message := ""

		// Validate email format
		if !isValidEmailFormat(email) {
			message = "Invalid email format."
		} else {
			domain := strings.Split(email, "@")[1]
			if !isValidDomain(domain) {
				message = "Invalid domain or no MX records found."
			} else if !checkSMTP(email) {
				message = "Email does not exist on the server."
			} else {
				message = "Email is valid and exists."
			}
		}

		// Send the result to the template
		tmpl.Execute(w, map[string]string{
			"Message": message,
		})
		return
	}
	tmpl.Execute(w, nil)
}

// Step 4: Web Template with CSS Styling
var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Email Verification Tool</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #1f1f1f;
            color: #f0f0f0;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .container {
            background-color: #2b2b2b;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0px 8px 16px rgba(0, 0, 0, 0.3);
            text-align: center;
            width: 350px;
        }
        h1 {
            font-size: 24px;
            color: #00d4ff;
            margin-bottom: 20px;
        }
        label {
            font-size: 16px;
            margin-bottom: 10px;
            display: block;
        }
        input[type="email"] {
            padding: 10px;
            width: 100%;
            border: none;
            border-radius: 4px;
            margin-bottom: 20px;
            box-sizing: border-box;
        }
        button {
            padding: 10px 20px;
            background-color: #00d4ff;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            color: #1f1f1f;
            font-size: 16px;
            width: 100%;
        }
        button:hover {
            background-color: #00a8cc;
        }
        p {
            margin-top: 20px;
            font-size: 14px;
            color: #ff7373;
        }
        .valid {
            color: #73ff73;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Email Verification Tool</h1>
        <form method="POST">
            <label for="email">Enter Email:</label>
            <input type="email" name="email" id="email" required>
            <button type="submit">Verify</button>
        </form>
        {{if .Message}}
        <p class="{{if eq .Message "Email is valid and exists."}}valid{{end}}">{{.Message}}</p>
        {{end}}
    </div>
</body>
</html>
`))

func main() {
	http.HandleFunc("/", emailCheckHandler)
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
