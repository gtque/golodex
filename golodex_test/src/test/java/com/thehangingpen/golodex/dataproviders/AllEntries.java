package com.thehangingpen.golodex.dataproviders;

import com.thegoate.data.DLProvider;
import com.thegoate.data.GoateDLP;
import com.thegoate.data.StaticDL;

@GoateDLP(name = "all entries")
public class AllEntries extends DLProvider {

	@Override
	public void init() {
		runData.put("run##", new StaticDL()
				.add("Scenario", "no sort, no search"))
			.put("run##", new StaticDL()
				.add("Scenario", "ascending sort, no search")
				.add("sort", "ascending"))
			.put("run##", new StaticDL()
				.add("Scenario", "descending sort, no search")
				.add("sort", "descending"))
			.put("run##", new StaticDL()
				.add("Scenario", "invalid sort, no search")
				.add("sort", "brownies"))
			.put("run##", new StaticDL()
				.add("Scenario", "no sort, simple search")
				.add("search", "ngel"));
	}
}
