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
