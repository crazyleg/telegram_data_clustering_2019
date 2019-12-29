package fasttext

// #cgo LDFLAGS: -L${SRCDIR} -lfasttext -lstdc++ -lm
// #include <stdlib.h>
// void load_model(char *name, char *path);
// int predict(char* name, char *query, float *prob, char **buf, int *count, int k, int buf_sz);
// int get_vector(char* name, char *query, float *results);
import "C"
import (
	"errors"
	"unsafe"
)

func LoadModel(name, path string) {
	p1 := C.CString(name)
	p2 := C.CString(path)

	C.load_model(p1, p2)

	C.free(unsafe.Pointer(p1))
	C.free(unsafe.Pointer(p2))
}

func GetVector(name, sentence string) ([55]float64, error) {
	np := C.CString(name)
	sentence += "\n"
	vector := make([]C.float, 55, 55)
	data := C.CString(sentence)

	C.get_vector(np, data, &vector[0])
	var result [55]float64;
	for i := range vector {
		result[i] = float64(vector[i])
	}
	return result, nil
}

// Predict - predict, return the topN predicted label and their corresponding probability
func Predict(name, sentence string, topN int) (map[string]float32, error) {
	result := make(map[string]float32)

	//add new line to sentence, due to the fasttext assumption
	sentence += "\n"

	cprob := make([]C.float, topN, topN)
	buf := make([]*C.char, topN, topN)
	var resultCnt C.int
	for i := 0; i < topN; i++ {
		buf[i] = (*C.char)(C.calloc(64, 1))
	}

	np := C.CString(name)
	data := C.CString(sentence)

	ret := C.predict(np, data, &cprob[0], &buf[0], &resultCnt, C.int(topN), 64)
	if ret != 0 {
		return result, errors.New("error in prediction")
	} else {
		for i := 0; i < int(resultCnt); i++ {
			result[C.GoString(buf[i])] = float32(cprob[i])
		}
	}
	//free the memory used by C
	C.free(unsafe.Pointer(data))
	C.free(unsafe.Pointer(np))
	for i := 0; i < topN; i++ {
		C.free(unsafe.Pointer(buf[i]))
	}

	return result, nil
}
