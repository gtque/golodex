package com.thehangingpen.golodex.test.core;

import com.thegoate.rest.RestCall;
import com.thegoate.staff.Employee;

public abstract class GolodexRestEmployee extends Employee {
	protected RestCall rest;

	@Override
	public String[] detailedScrub() {
		return new String[0];
	}

	@Override
	protected Employee init() {
		rest = new RestCall().header("X-GOLODEX-ID", getGolodexId())
			.header("X-GOLODEX-TOKEN", getGolodexToken());
		return this;
	}

	public GolodexRestEmployee setGolodexId(String id) {
		this.definition.put("golodex_id", id);
		return this;
	}

	public GolodexRestEmployee setGolodexToken(String token) {
		this.definition.put("golodex_token", token);
		return this;
	}

	public String getGolodexId() {
		return this.definition.get("golodex_id", "not_set", String.class);
	}

	public String getGolodexToken() {
		return this.definition.get("golodex_token", "not_set", String.class);
	}
}
