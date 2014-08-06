package main


func main() {
    c := NewContext(...) // Set up your context here
    http.Handle("/", ContextHandler{c, f1})
}