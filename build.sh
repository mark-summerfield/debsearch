cd cmd/cli
go build -o debsearch .
strip debsearch
mv debsearch ../..
cd ../..
