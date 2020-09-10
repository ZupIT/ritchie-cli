mod formula;

use std::env;

fn main() {
	let sample_text;
	let sample_list;
  let sample_bool;

	match env::var("SAMPLE_TEXT") {
		Ok(val) => sample_text = val,
		Err(_e) => sample_text = "none".to_string(),
  }

	match env::var("SAMPLE_LIST") {
		Ok(val) => sample_list = val,
		Err(_e) => sample_list = "none".to_string(),
  }

	match env::var("SAMPLE_BOOL") {
		Ok(val) => sample_bool = val,
		Err(_e) => sample_bool = "none".to_string(),
  }

	formula::run(sample_text, sample_list, sample_bool);
}
