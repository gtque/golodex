package com.thehangingpen.golodex.test.data;

import static com.thegoate.dsl.words.EutConfigDSL.eut;

import com.thehangingpen.golodex.test.core.GolodexRestEmployee;
import com.thegoate.staff.Employee;

public abstract class GolodexDataEmployee extends GolodexRestEmployee {

	@Override
	public Employee init(){
		super.init();
		rest.baseURL(eut("golodex.data.url", "https://golodex-data.thehangingpen.com"));
		return this;
	}
}
