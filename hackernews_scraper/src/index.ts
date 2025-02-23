// Initializes clients, assembles services, and runs the article workflow.

import config from './config/config';
import { createClients } from './clients/createClients';
import { createServices } from './services/createServices';
import { createWorkflow } from './workflow';

const EXIT_CODE_SUCCESS = 0;
const EXIT_CODE_FAILURE = 1;

const main = async (): Promise<number> => {
  // Initialize workflow and set default exit code
  let workflow: ReturnType<typeof createWorkflow> | null = null;
  let exitCode = EXIT_CODE_SUCCESS;
  try {
    // Create low-level infrastructure clients
    const clients = await createClients(config);

    // Assemble high-level services
    const services = createServices(clients);

    // Run the workflow
    workflow = createWorkflow(services);
    await workflow.run();

    console.info('Workflow completed successfully');
  } catch (error) {
    console.error('Workflow execution error:', error);
    exitCode = EXIT_CODE_FAILURE;
  } finally {
    // Shutdown workflow if initialized
    if (workflow) await workflow.shutdown();
  }
  return exitCode;
};

// Run main if this file is executed directly
if (require.main === module) {
  main().then((exitCode) => process.exit(exitCode));
}
