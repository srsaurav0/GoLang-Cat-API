go version

bee new cat-voting-app
cd cat-voting-app

cd GoLang-Cat-Api

python3 -m venv .venv
source .venv/bin/activate

go env GOPATH
export PATH=$PATH:$(go env GOPATH)/bin
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
bee version

go get github.com/beego/beego/v2
go mod tidy


npm init -y
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p


module.exports = {
  content: [
    "./views/**/*.{tpl,html}",   // Beego templates
    "./static/js/**/*.js"        // JavaScript files
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};


mkdir -p static/css
touch static/css/styles.css

@tailwind base;
@tailwind components;
@tailwind utilities;

npx tailwindcss -i ./static/css/styles.css -o ./static/css/output.css --watch


bee run


go test ./... -coverprofile=coverage.out

# Display total coverage percentage
go tool cover -func=coverage.out | grep total: | awk '{print $3}'

# Generate HTML coverage report (optional)
go tool cover -html=coverage.out -o coverage.html

# Open the HTML report (optional)
open coverage.html

https://github.com/srsaurav0/GoLang-Cat-API.git



python -m venv .venv
Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
.venv\Scripts\activate

go get github.com/beego/beego/v2@v2.3.4
go get github.com/beego/bee/v2@latest
bee version

bee run


Windows:
go test -coverprofile coverage.out ./...
go tool cover -html coverage.out