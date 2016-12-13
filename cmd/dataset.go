package cmd

type DataSet []Items

type Items []int

func OrderItems(d DataSet, h HeadTable) DataSet {
	// naive...?
	ds := make(DataSet, len(d))
	for x, dd := range d {
		items := []int{}
		for _, ih := range h {
			for _, id := range dd {
				if id == ih.Item {
					items = append(items, id)
				}
			}
		}
		ds[x] = items
	}
	return ds
}