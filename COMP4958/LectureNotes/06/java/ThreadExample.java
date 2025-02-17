class ThreadExample {
  public static void main(String... args) {
    Runnable r1 = () -> {
      while (true) {
        System.out.println("hello");
      }
    };

    Runnable r2 = () -> {
      while (true) {
        System.out.println("world");
      }
    };

    new Thread(r1).start();
    new Thread(r2).start();
  }
}
