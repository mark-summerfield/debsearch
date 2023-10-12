cd cmd/cli
go build -o debsearch .
strip debsearch
mv debsearch ../..
cd ../..

cd cmd/gui
go build -o DebFind .
strip DebFind
mv DebFind ../..
cd ../..
