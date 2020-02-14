package status_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	rabbitmqv1beta1 "github.com/pivotal/rabbitmq-for-kubernetes/api/v1beta1"
	rabbitmqstatus "github.com/pivotal/rabbitmq-for-kubernetes/internal/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("ClusterAvailable", func() {

	var (
		childServiceEndpoints *corev1.Endpoints
	)

	BeforeEach(func() {
		childServiceEndpoints = &corev1.Endpoints{}
	})

	When("at least one service endpoint is published", func() {
		var (
			conditionManager rabbitmqstatus.ClusterAvailableConditionManager
		)

		BeforeEach(func() {
			childServiceEndpoints.Subsets = []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{
						{
							IP: "1.2.3.4",
						},
						{
							IP: "5.6.7.8",
						},
					},
				},
			}
			conditionManager = rabbitmqstatus.NewClusterAvailableConditionManager(childServiceEndpoints)
		})

		It("returns the expected condition", func() {
			condition := conditionManager.Condition()
			By("having the correct type", func() {
				var conditionType rabbitmqv1beta1.RabbitmqClusterConditionType = "ClusterAvailable"
				Expect(condition.Type).To(Equal(conditionType))
			})

			By("having status true and reason message", func() {
				Expect(condition.Status).To(Equal(corev1.ConditionTrue))
				Expect(condition.Reason).To(Equal("AtLeastOneNodeAvailable"))
			})

			By("having a probe time", func() {
				Expect(condition.LastProbeTime).NotTo(Equal(metav1.Time{}))
			})
		})
	})

	When("no service endpoint is published", func() {
		var (
			conditionManager rabbitmqstatus.ClusterAvailableConditionManager
		)

		BeforeEach(func() {
			childServiceEndpoints.Subsets = []corev1.EndpointSubset{
				{
					Addresses: []corev1.EndpointAddress{},
				},
			}
			conditionManager = rabbitmqstatus.NewClusterAvailableConditionManager(childServiceEndpoints)
		})

		It("returns the expected condition", func() {
			condition := conditionManager.Condition()
			By("having the correct type", func() {
				var conditionType rabbitmqv1beta1.RabbitmqClusterConditionType = "ClusterAvailable"
				Expect(condition.Type).To(Equal(conditionType))
			})

			By("having status true and reason message", func() {
				Expect(condition.Status).To(Equal(corev1.ConditionFalse))
				Expect(condition.Reason).To(Equal("NoServiceEndpointsAvailable"))
				Expect(condition.Message).NotTo(BeEmpty())
			})

			By("having a probe time", func() {
				Expect(condition.LastProbeTime).NotTo(Equal(metav1.Time{}))
			})
		})
	})

	When("service endpoints do not exist", func() {
		var (
			conditionManager rabbitmqstatus.ClusterAvailableConditionManager
		)

		BeforeEach(func() {
			childServiceEndpoints = nil
			conditionManager = rabbitmqstatus.NewClusterAvailableConditionManager(childServiceEndpoints)
		})

		It("returns the expected condition", func() {
			condition := conditionManager.Condition()
			By("having the correct type", func() {
				var conditionType rabbitmqv1beta1.RabbitmqClusterConditionType = "ClusterAvailable"
				Expect(condition.Type).To(Equal(conditionType))
			})

			By("having status true and reason message", func() {
				Expect(condition.Status).To(Equal(corev1.ConditionFalse))
				Expect(condition.Reason).To(Equal("CouldNotAccessServiceEndpoints"))
				Expect(condition.Message).NotTo(BeEmpty())
			})

			By("having a probe time", func() {
				Expect(condition.LastProbeTime).NotTo(Equal(metav1.Time{}))
			})
		})
	})

})
