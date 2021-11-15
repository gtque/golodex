package com.thehangingpen.golodex.api;

import static com.thegoate.dsl.words.EutConfigDSL.eut;

import org.testng.annotations.Test;

import com.thegoate.Goate;
import com.thegoate.data.GoateProvider;
import com.thegoate.testng.TestNGEngineMethodDL;
import com.thehangingpen.golodex.test.api.all.GetAllApi;

public class PublicApiTests extends TestNGEngineMethodDL {

	@GoateProvider(name = "all entries")
	@Test(groups = {"api", "regression"}, dataProvider = "methodLoader")
	public void getAll(Goate testData) {
		Object result = new GetAllApi()
			.sort(get("sort", null, String.class))
			.search(get("search", null, String.class))
			.setGolodexId(get("golodex id", eut("golodex_id"), String.class))
			.setGolodexToken(get("golodex token", eut("golodex_token"), String.class))
			.build()
			.work();
	}

}
