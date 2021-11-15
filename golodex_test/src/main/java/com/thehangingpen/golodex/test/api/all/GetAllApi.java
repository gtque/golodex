package com.thehangingpen.golodex.test.api.all;

import com.thehangingpen.golodex.test.api.GolodexApiEmployee;

public class GetAllApi extends GolodexApiEmployee {

	public GetAllApi sort(String sort) {
		definition.put("sort", sort);
		return this;
	}

	public GetAllApi search(String search) {
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
