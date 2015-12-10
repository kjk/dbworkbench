@rem TODO: use godep

godep go vet github.com/kjk/dbworkbench

godep go build -o dbherohelper.exe
dbherohelper.exe -dev
