package com.github.alex_moon.noise;

import java.util.ArrayList;
import java.util.List;

import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

import org.json.simple.JSONObject;

import com.github.alex_moon.noise.core.Core;
import com.github.alex_moon.noise.fact.Fact;
import com.github.alex_moon.noise.fact.FactBean;
import com.github.alex_moon.noise.term.Term;

@Path("/")
public class Noise {
    @GET
    @Path("/{query}")
    @Produces(MediaType.APPLICATION_JSON)
    public List<FactBean> getFacts(@PathParam("query") String queryString) {
        Term query = Core.getTermController().getTerm(queryString, null);
        List<Fact> facts = Core.getFactController().getFactsForPrimaryTerm(query);
        List<FactBean> beans = new ArrayList<FactBean>();
        for (Fact fact : facts) {
            if (!fact.doubleValue().isNaN() && !fact.doubleValue().isInfinite()) {
                beans.add(new FactBean(fact));
            }
        }
        return beans;
    }
}
