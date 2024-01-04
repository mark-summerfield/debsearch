echo cmd/debsearch
cd cmd/debsearch
go build -o debsearch .
strip debsearch
mv debsearch ../..
cd ../..

echo cmd/DebFind
cd cmd/DebFind
go build -o DebFind .
strip DebFind
mv DebFind ../..
cd ../..
