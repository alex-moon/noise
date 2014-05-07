import java.io.IOException;
import java.io.PrintWriter;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

public class FuckServlet extends HttpServlet {
    public void doGet(HttpServletRequest req, HttpServletResponse res)
                    throws IOException, ServletException {

        res.setContentType("text/html");
        PrintWriter out = res.getWriter();
        out.println("");
        out.println("MyServlet");
        out.println("\t");
        out.println("");
        out.println("");
        out.println("Suck my peeny :)");
        out.println("");
        out.close();
    }
}
