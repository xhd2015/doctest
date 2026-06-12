Design an adversarial two-agent system.The main agent is responsible for:Analyzing and understanding the requirements
Writing tests, Sealing the doctests (to prevent them from being arbitrarily modified)

Once requirement confirmation is complete, the main agent enters the test design phase, ensuring the tests are comprehensive and complete. It then runs the tests to confirm that all test cases are in a failing state (RED).After that, the main agent invokes the sub-agent for implementation. It provides the sub-agent with the design document and an overview of the test cases.When the sub-agent needs clarification or confirmation on any issues during its work, it can call yield-pending-questions and then suspend the conversation.Upon receiving questions from the sub-agent, the main agent either resolves them independently or confirms them with the user, then feeds the answers back to the sub-agent.If the sub-agent reports that the work is complete, the main agent must:First verify that the staged tests have not been modified. If any modifications were made, they must have unavoidable/necessary justification.
Run the tests to ensure all results are correct (passing).

Sub-agent continuity is maintained via thread ID.

Main agent will be run as a skill, put the doc into /Users/xhd2015/Projects/xhd2015/agent-pro/agents/doctest/doc/DOC_STYLE_TEST_BASED_TDD.md

follow TDD style, write doctests first, run them and expect failure(for those expecting error, the error message would be wrong because the implementation may only be a stub return "error not implemented"), then run `git add <test-dir>` to seal tests, then goes to implementation, run the tests, fix code to address failures until all tests are green

The whole system is a practice of TDD: write tests that fail first, then correct implementation until success.