run: test

build:
	go build

test: build
	./cork-make -p "testfile" "testfile2" -r "date" -p "testfile3" -r "ls"

clean:
	rm -f cork-make
