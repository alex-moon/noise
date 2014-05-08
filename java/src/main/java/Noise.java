import java.io.IOException;
import java.io.PrintWriter;
import java.util.List;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.fact.Fact;
import com.github.alex_moon.noise.term.Term;

class Noise extends HttpServlet {
    public void doGet(HttpServletRequest req, HttpServletResponse res)
                    throws IOException, ServletException {

        String queryString = "car";
        Term query = Core.getTermController().getTerm(queryString, null);
        List<Fact> facts = Core.getFactController().getFactsForPrimaryTerm(query);
        
        res.setContentType("text/html");
        PrintWriter out = res.getWriter();
        out.println("We have a list of facts!");
        out.close();
    }
    
    public void init() {
        Core.getInstance().run();
    }
}
