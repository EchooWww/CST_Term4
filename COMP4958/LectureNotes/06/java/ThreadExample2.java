class ThreadExample2 implements Runnable {
  String  s;

  ThreadExample2(String s) {
    this.s = s;
  }

  public void run() {
    while (true) {
      System.out.println(s);
    }
  }

  public static void main(String... args) {
    new Thread(new ThreadExample2("hello")).start();
    new Thread(new ThreadExample2("world")).start();
  }
}
