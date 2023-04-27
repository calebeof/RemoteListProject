package remotelist

import (
	"errors"
	"fmt"
	"sync"
	"os"
	"io/ioutil"
	"strconv"
	"bufio"
)

type PersistentRemoteList struct {
	mu   sync.Mutex
	lists map[int][]int
}

type AppendArgs struct {
	ListIndex int
	Value int
}

const persistentDirectory = "./remoteLists/"

func (l *PersistentRemoteList) getListFromMemory(listIndex int) ([]int, error){
	if _, err := os.Stat(persistentDirectory); os.IsNotExist(err) {
		fmt.Printf("Persistent directory %s doesn't exists, creating now...\n", persistentDirectory)
		err = os.Mkdir(persistentDirectory, 0700)
		if err != nil {
			return nil, fmt.Errorf("Couldn't create persistent directory %s: %w\n", persistentDirectory, err)
		}
		return make([]int, 0), nil
	}
	fmt.Println("Directory exists.")

	files, err := ioutil.ReadDir(persistentDirectory)
    if err != nil {
        return nil, fmt.Errorf("Couldn't read persistent directory %s: %w\n", persistentDirectory, err)
    }
	fmt.Println("Read the directory.")

	for _, file := range files {
		fmt.Printf("filename: %s\n", file.Name())
		listIndexFile, err := strconv.Atoi(file.Name())
		if err != nil {
			fmt.Println("Couldn't convert listIndex file name...")
			continue
		}
		fmt.Println("Converted the filename.")
		if listIndexFile == listIndex {
			fmt.Printf("Found the list %d in memory!\n", listIndex)
			fileOpened, err := os.Open(persistentDirectory + file.Name())
			if err != nil {
				return nil, fmt.Errorf("Couldn't open file %s from list %d: %w.\n", file.Name(), listIndex, err)
			}
			defer fileOpened.Close()

			var list []int
			scanner := bufio.NewScanner(fileOpened)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				number, err := strconv.Atoi(scanner.Text())
				if err != nil {
					fmt.Println("Couldn't convert number in the list...")
					continue
				}
				fmt.Printf("Appending number %d to list %d...\n", number, listIndex)
				list = append(list, number)
			}

			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("Couldn't read numbers in file %s from list %d: %w.\n", file.Name(), listIndex, err)
			}
			
			return list, nil
		}
    }
	return nil, fmt.Errorf("Couldn't find the list in memory.\n")
}

func (l *PersistentRemoteList) saveListInMemory(listIndex int) (error){
	if _, err := os.Stat(persistentDirectory); os.IsNotExist(err) {
		fmt.Printf("Persistent directory %s doesn't exists, creating now...\n", persistentDirectory)
		err = os.Mkdir(persistentDirectory, 0700)
		if err != nil {
			return fmt.Errorf("Couldn't create persistent directory %s: %w\n", persistentDirectory, err)
		}
	}

	fileName := persistentDirectory + strconv.Itoa(listIndex)
	file, err := os.Create(fileName)
	defer file.Close()

    if err != nil {
        return fmt.Errorf("Couldn't create persistent file %s: %w.\n", fileName, err)
    }

	for _, number := range l.lists[listIndex] {
		writeNumber := fmt.Sprintf("%s\n", strconv.Itoa(number))
		_, err = file.WriteString(writeNumber)
		if err != nil {
			return fmt.Errorf("Couldn't save number %d in file %s: %w.\n", number, fileName, err)
		}
	}

	fmt.Printf("Successfully saved list %d in memory at %s.\n", listIndex, fileName)
	return nil
}

func (l *PersistentRemoteList) Append(args AppendArgs, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, found := l.lists[args.ListIndex]
	if !found {
		fmt.Printf("List %d doesn't exists in runtime, let's see if it's in memory...\n", args.ListIndex)
		list, err := l.getListFromMemory(args.ListIndex)
		if err != nil {
			fmt.Printf("List isn't persisted, so we'll create it...\n")
		}else {
			fmt.Printf("List is persisted, now recovering it...\n")
			l.lists[args.ListIndex] = list
			fmt.Printf("List is %v\n", list)
		}
	}

	l.lists[args.ListIndex] = append(l.lists[args.ListIndex], args.Value)
	fmt.Println(l.lists[args.ListIndex])
	err := l.saveListInMemory(args.ListIndex)
	if err != nil {
		fmt.Printf("Couldn't save list %d in memory: %s.\n", args.ListIndex, err)
	}
	*reply = true
	return nil
}

func (l *PersistentRemoteList) Remove(listIndex int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	list, found := l.lists[listIndex]
	if !found {
		return fmt.Errorf("List %d doesn't exists.\n", listIndex)
	}

	if len(list) > 0 {
		*reply = list[len(list)-1]
		l.lists[listIndex] = list[:len(list)-1]
	} else {
		return errors.New("empty list")
	}
	return nil
}

func (l *PersistentRemoteList) Get(listIndex int, reply *[]int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	list, found := l.lists[listIndex]
	if !found {
		fmt.Printf("List %d doesn't exists in runtime, let's see if it's in memory...\n", listIndex)
		list, err := l.getListFromMemory(listIndex)
		if err != nil {
			return fmt.Errorf("List %d doesn't exists.\n", listIndex)
		}else {
			fmt.Printf("List is persisted, now recovering it...\n")
			l.lists[listIndex] = list
			fmt.Printf("List is %v\n", list)
		}
	}
	*reply = list
	return nil
}

func (l *PersistentRemoteList) Size(listIndex int, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	list, found := l.lists[listIndex]
	if !found {
		fmt.Printf("List %d doesn't exists in runtime, let's see if it's in memory...\n", listIndex)
		list, err := l.getListFromMemory(listIndex)
		if err != nil {
			return fmt.Errorf("List %d doesn't exists.\n", listIndex)
		}else {
			fmt.Printf("List is persisted, now recovering it...\n")
			l.lists[listIndex] = list
			fmt.Printf("List is %v\n", list)
		}
	}
	*reply = len(list)
	return nil
}

func NewPersistentRemoteList() *PersistentRemoteList {
	remoteList := new(PersistentRemoteList)
	remoteList.lists = make(map[int][]int)
	return remoteList
}
