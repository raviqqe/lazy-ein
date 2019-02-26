define i8* @atomic.load(i8**) {
  %2 = load atomic i8*, i8** %0 seq_cst, align 8
  ret i8* %2
}

define void @atomic.store(i8*, i8**) {
  store atomic i8* %0, i8** %1 seq_cst, align 8
  ret void
}

define i1 @atomic.cmpxchg(i8**, i8*, i8*) {
  %4 = cmpxchg i8** %0, i8* %1, i8* %2 seq_cst seq_cst
  %5 = extractvalue { i8*, i1 } %4, 1
  ret i1 %5
}
