package com.ritchie;

import com.ritchie.formula.Hello;

public class Main {

    public static void main(String[] args) {
        String input1 = System.getenv("SAMPLE_TEXT");
        String input2 = System.getenv("SAMPLE_LIST");
        boolean input3 = Boolean.parseBoolean(System.getenv("SAMPLE_BOOL"));
        Hello hello = new Hello(input1, input2, input3);
        System.out.println(hello.Run());
    }
}