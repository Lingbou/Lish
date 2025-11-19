package parser

// Pipeline 表示一个管道
type Pipeline struct {
	Commands []*ParsedCommand
}

// Statement 表示一条完整语句
type Statement struct {
	Pipelines []*Pipeline
	Operator  TokenType // And, Or, Semicolon
}

// ParsePipeline 解析管道命令
func ParsePipeline(input string) (*Statement, error) {
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()

	statement := &Statement{
		Pipelines: make([]*Pipeline, 0),
	}

	currentPipeline := &Pipeline{
		Commands: make([]*ParsedCommand, 0),
	}

	currentCmd := &ParsedCommand{
		Args: make([]string, 0),
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token.Type {
		case TokenWord:
			if currentCmd.Command == "" {
				currentCmd.Command = token.Value
			} else {
				currentCmd.Args = append(currentCmd.Args, token.Value)
			}

		case TokenPipe:
			// 完成当前命令，添加到管道
			if currentCmd.Command != "" {
				currentPipeline.Commands = append(currentPipeline.Commands, currentCmd)
				currentCmd = &ParsedCommand{Args: make([]string, 0)}
			}

		case TokenRedirectOut, TokenRedirectAppend:
			// 读取下一个 token 作为文件名
			if i+1 < len(tokens) && tokens[i+1].Type == TokenWord {
				currentCmd.RedirectOut = tokens[i+1].Value
				currentCmd.Append = (token.Type == TokenRedirectAppend)
				i++ // 跳过文件名 token
			}

		case TokenRedirectIn:
			if i+1 < len(tokens) && tokens[i+1].Type == TokenWord {
				currentCmd.RedirectIn = tokens[i+1].Value
				i++
			}

		case TokenRedirectErr:
			if i+1 < len(tokens) && tokens[i+1].Type == TokenWord {
				currentCmd.RedirectErr = tokens[i+1].Value
				i++
			}

		case TokenAnd, TokenOr, TokenSemicolon:
			// 完成当前命令和管道
			if currentCmd.Command != "" {
				currentPipeline.Commands = append(currentPipeline.Commands, currentCmd)
				currentCmd = &ParsedCommand{Args: make([]string, 0)}
			}
			if len(currentPipeline.Commands) > 0 {
				statement.Pipelines = append(statement.Pipelines, currentPipeline)
				statement.Operator = token.Type
				currentPipeline = &Pipeline{Commands: make([]*ParsedCommand, 0)}
			}

		case TokenEOF:
			// 完成最后的命令和管道
			if currentCmd.Command != "" {
				currentPipeline.Commands = append(currentPipeline.Commands, currentCmd)
			}
			if len(currentPipeline.Commands) > 0 {
				statement.Pipelines = append(statement.Pipelines, currentPipeline)
			}
		}
	}

	return statement, nil
}
