package tpl_java

const (
	TemplateMain = `import {{bin-name}}.{{bin-name-first-upper}};

public class Main {

    public static void main(String[] args) throws Exception {
        String input1 = System.getenv("SAMPLE_TEXT");
        String input2 = System.getenv("SAMPLE_LIST");
        boolean input3 = Boolean.parseBoolean(System.getenv("SAMPLE_BOOL"));
        {{bin-name-first-upper}} {{bin-name}} = new {{bin-name-first-upper}}(input1, input2, input3);
        {{bin-name}}.Run();
    }
}`

	TemplateMakefile = `# Java Parameters
BINARY_NAME_UNIX={{bin-name}}.sh
BINARY_NAME_WINDOWS={{bin-name}}.bat
DIST=../dist
DIST_DIR=$(DIST)/commons/bin

build:
	mkdir -p $(DIST_DIR)
	javac -source 1.8 -target 1.8 *.java
	echo "Main-Class: Main" > manifest.txt
	jar cvfm Main.jar manifest.txt *.class {{bin-name}}/*.class
	cp run_template $(BINARY_NAME_UNIX) && chmod +x $(BINARY_NAME_UNIX)
	cp run_template $(BINARY_NAME_WINDOWS) && chmod +x $(BINARY_NAME_WINDOWS)
	cp Main.jar $(BINARY_NAME_WINDOWS) $(BINARY_NAME_UNIX) $(DIST_DIR)
	#Clean files
	rm Main.jar manifest.txt *.class {{bin-name}}/*.class $(BINARY_NAME_WINDOWS) $(BINARY_NAME_UNIX)`

	TemplateRunTemplate = `#!/bin/sh
java -jar Main.jar`

	TemplateFileJava = `package {{bin-name}};

public class {{bin-name-first-upper}} {

    private String input1;
    private String input2;
    private boolean input3;

    public void Run() throws Exception {
        System.out.printf("Hello World!\n");
        System.out.printf("You receive %s in text.\n", input1);
        System.out.printf("You receive %s in list.\n", input2);
        System.out.printf("You receive %s in boolean.\n", input3);
    }

    public {{bin-name-first-upper}}(String input1, String input2, boolean input3) {
        this.input1 = input1;
        this.input2 = input2;
        this.input3 = input3;
    }

    public String getInput1() {
        return input1;
    }

    public void setInput1(String input1) {
        this.input1 = input1;
    }

    public String getInput2() {
        return input2;
    }

    public void setInput2(String input2) {
        this.input2 = input2;
    }

    public boolean isInput3() {
        return input3;
    }

    public void setInput3(boolean input3) {
        this.input3 = input3;
    }
}`
)
