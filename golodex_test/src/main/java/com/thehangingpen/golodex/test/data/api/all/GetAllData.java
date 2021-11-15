package com.thehangingpen.golodex.test.data.api.all;

import com.thehangingpen.golodex.test.data.GolodexDataEmployee;

public class GetAllData extends GolodexDataEmployee {

	public GetAllData sort(String sort) {
		definition.put("sort", sort);
		return this;
	}

	public GetAllData search(String search) {
		definition.put("search", search);
		return this;
	}

	@Override
	protected Object doWork() {
		if(definition.get("sort", null) != null){
			rest.queryParam("sort", definition.get("sort"));
		}
		if(definition.get("search", null) != null){
			rest.queryParam("search", definition.get("search"));
		}
		return rest.get("api/all");
	}
}
