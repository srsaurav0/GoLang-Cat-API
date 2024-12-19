go version

bee new cat-voting-app
cd cat-voting-app

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