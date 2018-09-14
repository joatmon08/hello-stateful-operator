package hellostateful

func labelsForHelloStateful(name string) map[string]string {
	return map[string]string{"service": "app", "cr": name}
}

func labelsForHelloStatefulBackup(name string) map[string]string {
	return map[string]string{"service": "backup", "cr": name}
}
