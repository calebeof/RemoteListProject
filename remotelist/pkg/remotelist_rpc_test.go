package remotelist

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestGetListFromMemory(t *testing.T) {
	testCases := []struct {
		name      string
		l         *PersistentRemoteList
		listIndex int
		want      []int
		wantErr   bool
	}{
		{
			name:      "non-existent list in memory",
			l:         &PersistentRemoteList{},
			listIndex: 0,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "non-existent directory",
			l:         &PersistentRemoteList{},
			listIndex: 0,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "directory exists but file doesn't",
			l:         &PersistentRemoteList{},
			listIndex: 1,
			want:      nil,
			wantErr:   true,
		},
		{
			name: "list exists in memory",
			l: &PersistentRemoteList{
				lists: map[int][]int{
					2: {1, 2, 3},
				},
			},
			listIndex: 2,
			want:      []int{1, 2, 3},
			wantErr:   false,
		},
		{
			name:      "file exists and list is recovered",
			l:         &PersistentRemoteList{},
			listIndex: 3,
			want:      []int{1, 2, 3},
			wantErr:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := os.Stat(persistentDirectory); os.IsNotExist(err) {
				os.Mkdir(persistentDirectory, 0700)
			}

			if testCase.name == "file exists and list is recovered" || testCase.name == "list exists in memory" {
				file, err := os.Create(persistentDirectory + strconv.Itoa(testCase.listIndex))
				if err != nil {
					t.Fatalf("Error creating file: %v", err)
				}
				defer file.Close()
				file.WriteString("1\n2\n3\n")
			}

			defer func() {
				err := os.RemoveAll(persistentDirectory)
				if err != nil {
					t.Fatalf("failed to remove test directory: %v", err)
				}
			}()

			got, err := testCase.l.getListFromMemory(testCase.listIndex)
			if (err != nil) != testCase.wantErr {
				t.Errorf("getListFromMemory() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("getListFromMemory() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestSaveListInMemory(t *testing.T) {
	testCases := []struct {
		name           string
		listIndex      int
		l              *PersistentRemoteList
		expectedResult error
	}{
		{
			name:           "non-existent directory",
			listIndex:      1,
			l:              &PersistentRemoteList{},
			expectedResult: nil,
		},
		{
			name:           "existent directory",
			listIndex:      2,
			l:              &PersistentRemoteList{},
			expectedResult: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.l.saveListInMemory(testCase.listIndex)
			if err != testCase.expectedResult {
				t.Errorf("Expected %v but got %v", testCase.expectedResult, err)
			}
			defer func() {
				err := os.RemoveAll(persistentDirectory)
				if err != nil {
					t.Fatalf("failed to remove test directory: %v", err)
				}
			}()
		})
	}
}

func TestAppend(t *testing.T) {
	l := &PersistentRemoteList{
		lists: make(map[int][]int),
	}

	err := os.MkdirAll(persistentDirectory, 0700)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	defer func() {
		err := os.RemoveAll(persistentDirectory)
		if err != nil {
			t.Fatalf("failed to remove test directory: %v", err)
		}
	}()

	testCases := []struct {
		listIndex   int
		value       int
		expectedLen int
	}{
		{
			listIndex:   0,
			value:       1,
			expectedLen: 1,
		},
		{
			listIndex:   0,
			value:       2,
			expectedLen: 2,
		},
		{
			listIndex:   1,
			value:       3,
			expectedLen: 1,
		},
		{
			listIndex:   1,
			value:       4,
			expectedLen: 2,
		},
	}

	// run the test cases
	for _, testCase := range testCases {
		err := l.Append(AppendArgs{ListIndex: testCase.listIndex, Value: testCase.value}, new(bool))
		if err != nil {
			t.Errorf("Append(%v) failed with error: %v", testCase, err)
		}
		if len(l.lists[testCase.listIndex]) != testCase.expectedLen {
			t.Errorf("Append(%v) returned list length %v, expected %v", testCase, len(l.lists[testCase.listIndex]), testCase.expectedLen)
		}
	}

	// check that the persistent data was saved correctly
	for i := 0; i < 2; i++ {
		filePath := persistentDirectory + strconv.Itoa(i)
		file, err := os.Open(filePath)
		if err != nil {
			t.Errorf("failed to open file %v: %v", filePath, err)
			continue
		}
		defer file.Close()

		var list []int
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			number, err := strconv.Atoi(scanner.Text())
			if err != nil {
				t.Errorf("failed to convert number in file %v: %v", filePath, err)
				continue
			}
			list = append(list, number)
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("failed to read numbers in file %v: %v", filePath, err)
		}
		if len(list) != len(l.lists[i]) {
			t.Errorf("list %v saved in file %v has length %v, expected %v", i, filePath, len(list), len(l.lists[i]))
		}
		for j := 0; j < len(list); j++ {
			if list[j] != l.lists[i][j] {
				t.Errorf("list %v saved in file %v is different from the in-memory list", i, filePath)
			}
		}
	}

	defer func() {
		err := os.RemoveAll(persistentDirectory)
		if err != nil {
			t.Fatalf("failed to remove test directory: %v", err)
		}
	}()
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name          string
		args          []AppendArgs
		index         int
		expectedList  []int
		expectedError error
	}{
		{
			name: "Get empty list",
			args: []AppendArgs{
				//{ListIndex: 0, Value: 1},
			},
			index:         0,
			expectedList:  []int{},
			expectedError: nil,
		},
		{
			name: "Get non-empty list",
			args: []AppendArgs{
				{ListIndex: 0, Value: 1},
				{ListIndex: 0, Value: 2},
				{ListIndex: 0, Value: 3},
			},
			index:         0,
			expectedList:  []int{1, 2, 3},
			expectedError: nil,
		},
	}

	l := &PersistentRemoteList{lists: make(map[int][]int)}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Append values to the list
			for _, args := range testCase.args {
				err := l.Append(args, new(bool))
				if err != nil {
					t.Errorf("Unexpected error on Append: %v", err)
				}
			}

			defer func() {
				err := os.RemoveAll(persistentDirectory)
				if err != nil {
					t.Fatalf("failed to remove test directory: %v", err)
				}
			}()

			// Call the Get method and check its result
			var got []int
			err := l.Get(testCase.index, &got)
			if len(got) != len(testCase.expectedList) {
				t.Errorf("Unexpected list: got %v, expected %v", got, testCase.expectedList)
			}

			for i := 0; i < len(got); i++ {
				if got[i] != testCase.expectedList[i] {
					t.Errorf("Unexpected list: got %v, expected %v", got, testCase.expectedList)
				}
			}

			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("Unexpected error: got %v, expected %v", err, testCase.expectedError)
			}
		})
	}
}
