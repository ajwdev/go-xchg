package main

import (
	"fmt"
	"github.com/williamsandrew/go-xchg"
	"log"
	"os"
	"syscall"
)

func getInode(path string) (uint64, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	return fi.Sys().(*syscall.Stat_t).Ino, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("must provide source and destination files")
	}

	// Print out the inode numbers pre exchange
	inode1, err := getInode(os.Args[1])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	inode2, err := getInode(os.Args[2])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	fmt.Printf("%s inode: %d\n", os.Args[1], inode1)
	fmt.Printf("%s inode: %d\n", os.Args[2], inode2)

	// Do the exchange
	fmt.Println("exchanging...")
	err = xchg.Exchange(os.Args[1], os.Args[2])
	if err != nil {
		if err == syscall.ENOSYS {
			log.Fatalf("system does not support 'renameat2' system call. error: %s\n", err)
		}

		log.Fatalf("error: %s\n", err)
	}

	// Print out the new inode numbers
	inode1, err = getInode(os.Args[1])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	inode2, err = getInode(os.Args[2])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	fmt.Printf("%s inode: %d\n", os.Args[1], inode1)
	fmt.Printf("%s inode: %d\n", os.Args[2], inode2)
}
