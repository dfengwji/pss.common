package proxy

const (
	/**
	管理员数据表
	*/
	TableAdmin = "admins"
	/**
	用户数据表
	*/
	TableUser = "users"
	/**
	智能笔硬件数据表
	*/
	TablePen = "pens"
	/**
	微信用户数据表
	*/
	TableWechat = "wechat"
	/**
	上传的文件数据表
	*/
	TableAsset = "assets"

	/**
	实时保存从客户端收到的每条笔迹碎片，做标记，定时删除
	*/
	TableFragment = "fragments"

	/**
	针对一道试题的完整的笔迹数据表
	*/
	TableWriting = "writing"

	/**
	试题内容答案表，题库最主要的表
	*/
	TableExercises = "exercises"

	/**
	试题标签表
	*/
	TableExercisesTag = "examTags"

	/**
	题型表
	*/
	TableExamTypes = "examTypes"

	/**
	试题样式表
	*/
	TableOriginStyle = "originStyles"

	/**
	试题的来源:书或者试卷表
	*/
	TableOriginal = "originBooks"

	/**
	组卷生成的新的练习题册或者试卷表
	*/
	TablePublicBook = "publicBooks"

	/**
	组卷生成新的作业本样式表
	*/
	TablePublicStyle = "publicStyles"

	/**
	组卷生成新的错题本样式表
	*/
	TablePrivateStyle = "privateStyles"

	/**
	错题本
	*/
	TablePrivateBook = "privateBooks"

	/**
	阅卷结果表
	*/
	TableMarking = "marking"

	/**
	数量自动增加
	*/
	TableSequence = "sequences"

	/**
	用户地址表
	*/
	TableAddress = "address"

	/**
	用户答题计数关系表
	*/
	TableCount = "counts"

	/**
	书页表
	*/
	TablePages = "pages"

	TableScene = "scenes"

	TableClasses = "classes"

	TableStudent = "students"

	TableTeacher = "teachers"

	/*
		学生申请加入机构
	*/
	TableApply = "applies"

	TableCourseMenu = "courseMenus"
	TableMicroCourse = "courses"
	TableCourseDraft = "drafts"
	TableDraftWriting = "draftWritings"
	TableNoteBook = "notebooks"
	TableNoteWriting = "noteWritings"
	TableRoleID = "roleIDs"
	TableBookID = "bookIDs"
)
