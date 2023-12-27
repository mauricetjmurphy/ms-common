package criteria

func InUint(column string, values ...uint) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InUint32(column string, values ...uint32) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InUint64(column string, values ...uint64) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InInt(column string, values ...int) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InInt32(column string, values ...int) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InInt64(column string, values ...int64) Expr {
	if len(values) > 0 {
		return Eq(column, values)
	}
	return nil
}

func InStr(column string, values ...string) Expr {
	if len(values) == 0 {
		return nil
	}
	return Eq(column, values)
}

func NullOrInt64(column string, values []int64) Expr {
	if len(values) == 0 {
		return IsNil(column)
	}
	return InInt64(column, values...)
}

func NullOrUint(column string, values []uint) Expr {
	if len(values) == 0 {
		return IsNil(column)
	}
	return InUint(column, values...)
}
