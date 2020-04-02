resource "aws_cloudwatch_event_rule" "gas_express_updater_cloudwatch_rule" {
  name                = "gas-express-updater"
  description         = "Cloudwatch rule to trigger gas express updater task."
  schedule_expression = var.schedule
}

resource "aws_cloudwatch_event_target" "gas_express_updater_cloudwatch_target" {
  arn  = aws_lambda_function.gas-express-updater-lambda.arn
  rule = aws_cloudwatch_event_rule.gas_express_updater_cloudwatch_rule.name
}

resource "aws_lambda_permission" "allow_cloudwatch_trigger_lambda" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.gas-express-updater-lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.gas_express_updater_cloudwatch_rule.arn
}
