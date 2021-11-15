package com.thehangingpen.golodex.test.api;

import static com.thegoate.dsl.words.EutConfigDSL.eut;

import com.thehangingpen.golodex.test.core.GolodexRestEmployee;
import com.thegoate.staff.Employee;

public abstract class GolodexApiEmployee extends GolodexRestEmployee {

	@Override
	public Employee init(){
		super.init();
		rest.baseURL(eut("golodex.api.url", "https://golodex.thehangingpen.com"));
		return this;
	}
}
