interface LogData {
  request_id: string;
  method: string;
  path: string;
  query_params: Record<string, string>;
  status_code: number;
  duration_ms: number;
  timestamp: string;
  client_ip: string;
  user_agent: string;
  content_type: string;
  request_size_bytes: number;
  response_size_bytes: number;
  locale?: string;
  timezone?: string;
  service_name: string;
  service_version: string;
  git_sha?: string;
  host: string;
  environment: string;
  error_type?: string;
  error_message?: string;
  error_stack?: string;
  referrer?: string;
  accept_language?: string;
  memory_used_mb?: number;
}

function generateRequestId(): string {
  return `req_${Math.random().toString(36).substring(2, 11)}`;
}

function getClientIp(req: Request): string {
  const forwarded = req.headers.get("x-forwarded-for");
  const realIp = req.headers.get("x-real-ip");
  const cfConnectingIp = req.headers.get("cf-connecting-ip");

  return cfConnectingIp || realIp || forwarded?.split(",")[0] || "unknown";
}

function parseQueryParams(url: URL): Record<string, string> {
  const params: Record<string, string> = {};
  url.searchParams.forEach((value, key) => {
    params[key] = value;
  });
  return params;
}

export function createLogger() {
  const serviceName = process.env.SERVICE_NAME || "blog";
  const serviceVersion = process.env.SERVICE_VERSION || "1.0.0";
  const gitSha = process.env.GIT_SHA;
  const environment = process.env.NODE_ENV || "development";
  const host = process.env.HOSTNAME || Bun.hostname;

  return async function log(
    req: Request,
    res: Response,
    startTime: number,
    error?: Error
  ) {
    const endTime = performance.now();
    const url = new URL(req.url);
    const requestId = generateRequestId();

    const logData: LogData = {
      request_id: requestId,
      method: req.method,
      path: url.pathname,
      query_params: parseQueryParams(url),
      status_code: res.status,
      duration_ms: Math.round(endTime - startTime),
      timestamp: new Date().toISOString(),
      client_ip: getClientIp(req),
      user_agent: req.headers.get("user-agent") || "unknown",
      content_type: res.headers.get("content-type") || "unknown",
      request_size_bytes: parseInt(req.headers.get("content-length") || "0"),
      response_size_bytes: parseInt(res.headers.get("content-length") || "0"),
      service_name: serviceName,
      service_version: serviceVersion,
      host,
      environment,
      referrer: req.headers.get("referer") || undefined,
      accept_language: req.headers.get("accept-language") || undefined,
      git_sha: gitSha,
      memory_used_mb: Math.round(process.memoryUsage().heapUsed / 1024 / 1024),
    };

    // Add locale/timezone from Accept-Language header
    const acceptLang = req.headers.get("accept-language");
    if (acceptLang) {
      const locale = acceptLang.split(",")[0].split(";")[0].trim();
      logData.locale = locale;
    }

    // Add error details if present
    if (error) {
      logData.error_type = error.name;
      logData.error_message = error.message;
      logData.error_stack = error.stack;
    }

    // Log as JSON
    console.log(JSON.stringify(logData));
  };
}
