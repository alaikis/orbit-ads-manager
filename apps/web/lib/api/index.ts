// API Client
export { apiClient, ApiError } from './client'

// Types
export type {
  User,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  RefreshRequest,
  Store,
  CreateStoreRequest,
  UpdateStoreRequest,
  Campaign,
  CampaignMetrics,
  AdGroup,
  MetricsSummary,
  Notification,
  AdAccount,
  Product,
  Feed,
  Report,
  Rule,
  CreateRuleRequest,
  Workspace,
  Provider,
  PaginatedResponse,
  ApiResponse,
} from './types'

// Services
export { authService } from './services/auth.service'
export { metricsService } from './services/metrics.service'
export { notificationService } from './services/notification.service'
export { productService } from './services/product.service'
export { feedService } from './services/feed.service'
export { reportService } from './services/report.service'
export { ruleService } from './services/rule.service'
export { workspaceService } from './services/workspace.service'
export { connectionService } from './services/connection.service'
export { providerService } from './services/provider.service'
export { adAccountService } from './services/ad-account.service'

// Service classes (for inheritance/extension)
export { AuthService } from './services/auth.service'
export { MetricsService } from './services/metrics.service'
export { NotificationService } from './services/notification.service'
export { ProductService } from './services/product.service'
export { FeedService } from './services/feed.service'
export { ReportService } from './services/report.service'
export { RuleService } from './services/rule.service'
export { WorkspaceService } from './services/workspace.service'
export { ProviderService } from './services/provider.service'

// Legacy api object for backward compatibility
import { apiClient } from './client'
export const api = {
  get: <T>(path: string) => apiClient.get<T>(path),
  post: <T>(path: string, data?: unknown) => apiClient.post<T>(path, data),
  patch: <T>(path: string, data?: unknown) => apiClient.patch<T>(path, data),
  delete: <T>(path: string) => apiClient.delete<T>(path),
}
