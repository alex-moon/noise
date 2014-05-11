package com.github.alex_moon.noise;

import java.io.IOException;
import java.io.PrintWriter;
import java.util.ArrayList;
import java.util.List;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.json.simple.JSONObject;
import org.json.simple.JSONValue;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.fact.Fact;
import com.github.alex_moon.noise.term.Term;

public class Noise extends HttpServlet {
    public void doGet(HttpServletRequest req, HttpServletResponse res)
                    throws IOException, ServletException {
        res.setContentType("text/html");
        PrintWriter out = res.getWriter();

        String queryString = req.getParameter("q");
        if (queryString == null) {
            out.println("<form method='GET' action='.'><input type='text' value='search' name='q' /><input type='submit' value='Search!' /></form>");
        } else {
            Term query = Core.getTermController().getTerm(queryString, null);
            List<Fact> facts = Core.getFactController().getFactsForPrimaryTerm(query);
            List<JSONObject> factsToJson = new ArrayList<JSONObject>();
            for (Fact fact: facts) {
                if (!fact.doubleValue().isNaN() && !fact.doubleValue().isInfinite()) {
                    factsToJson.add(fact.toJson());
                }
            }
            out.println(JSONValue.toJSONString(factsToJson));
        }
        out.close();
    }
}
