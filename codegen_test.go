package codegen_test

import (
	"fmt"
	"github.com/u-mulder/codegen"
)

//func ExampleFoo()     // documents the Foo function or type
//func ExampleBar_Qux() // documents the Qux method of type Bar
//func Example()        // documents the package as a whole

func ExampleAddGetExistingSnippet() {
	c, _ := codegen.New()

	c.AddSnippet("SnippetNo1", "SomeCodeGoesHere")

	fmt.Println(c.GetSnippet("SnippetNo1"))
	// Output: SomeCodeGoesHere <nil>
}

func ExampleAddGetNonExistingSnippet() {
	c, _ := codegen.New()

	c.AddSnippet("SnippetNo1", "SomeCodeGoesHere")

	fmt.Println(c.GetSnippet("SnippetNo2"))
	// Output: Snippet with such key not found
}

func ExampleEncloseInSingleQuotes() {
	fmt.Println(codegen.EncloseInSingleQuotes("VALUE"))
	// Output: 'VALUE'
}

func ExampleCodegen_AddDefaultSnippets() {
	c, _ := codegen.New()

	c.AddDefaultSnippets()

	fmt.Println(c.GetSnippet("mevent_obj"))
	// Output: $meo = new CEventType; <nil>
}

func ExampleCodegen_Generate() {
	c, _ := codegen.New()

	c.AddSnippet("php_quick_echo", "<?=")
	c.AddSnippet("php_closing_tag", "?>")

	c.RegisterGenerator("qe_script", quickEcho)
	r, _ := c.Generate("qe_script")

	fmt.Println(r)
	// Output: <?='123-321'?>
}

func ExampleCodegen_RegisterDefaultGenerators() {
	c, _ := codegen.New()

	c.RegisterDefaultGenerators()
	_, errOne := c.Generate("uf")
	_, errTwo := c.Generate("ibprop")
	_, errThree := c.Generate("non_ex_g")

	fmt.Println(errOne, errTwo, errThree)
	// Output: <nil> <nil> Generator with such name not found
}

func quickEcho(c *codegen.Codegen) string {
	var r, tmpStr string

	tmpStr, _ = c.GetSnippet("php_quick_echo")
	if "" != tmpStr {
		r = tmpStr
	}

	r += codegen.EncloseInSingleQuotes("123-321")

	tmpStr, _ = c.GetSnippet("php_closing_tag")
	if "" != tmpStr {
		r += tmpStr
	}

	return r
}
