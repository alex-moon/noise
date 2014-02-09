import com.github.alex_moon.noise.text.Controller;
import com.github.alex_moon.noise.text.Text;

class HelloWorldApp {
    public static void main(String[] args) {
        String textString = "";
        Controller textController = new Controller();
        try {
            textString = textController.getOneText().toString();
        } catch (Exception e) {
            System.out.println("Ah, the ol catch-all eh? Works every time");
            e.printStackTrace();
        }
        System.out.println(textString); // Display the string.
    }
}
